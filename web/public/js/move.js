import Vector, { Vector4 } from '/js/vector.js';

const SNAP_THRESHOLD = 20;

const Directions = {
    UP: 0,
    RIGHT: 1,
    DOWN: 2,
    LEFT: 3,

    fromRect(dir) {
        if (dir == "y") {
            return Directions.UP;
        } else if (dir == "w") {
            return Directions.RIGHT;
        } else if (dir == "z") {
            return Directions.DOWN;
        }

        return Directions.LEFT;
    }
};

class ScreenMover {
    constructor(alpine) {
        this.alpine = alpine;

        this.canvas = document.querySelector("#screens");

        this.target = null;
        this.startPos = null;
        this.currPos = null;

        document.onmouseup = this._handleDragEnd.bind(this);
    }

    /**
     * Center the canvas in the window based on its contents size
     */
    centerCanvas() {
        let width = 0;
        let height = 0;

        for (let g = 0; g < this.alpine.groups.length; g++) {
            let group = this.alpine.groups[g];       

            for (let s = 0; s < group.screens.length; s++) {
                let screen = group.screens[s];
                width = Math.max(width, group.pos.x + screen.pos.x + screen.width);
                height = Math.max(height, group.pos.y + screen.pos.y + screen.height);
            }
        }

        this.alpine.canvas = { width, height };
    }

    /**
     * event handler for clicking on a screen group
     */
    handleDragStart(e, group) {
        this.target = group;
        this.startPos = new Vector(e.clientX, e.clientY);
        this.currPos = new Vector(e.clientX, e.clientY);

        document.onmousemove = this._handleDragMove.bind(this);
    }

    /**
     * Get an array of screens that this group overlaps
     */
    findOverlappingScreens(group) {
        let overlapping = [];
        let groupEdges = edges(group, { pos: Vector.zero });

        for (let g = 0; g < this.alpine.groups.length; g++) {
            if (this.alpine.groups[g].id == group.id) {
                continue;
            }

            for (let s = 0; s < this.alpine.groups[g].screens.length; s++) {
                let screen = this.alpine.groups[g].screens[s];
                let screenEdges = edges(screen, { pos: Vector.zero }, this.alpine.groups[g]);

                if (groupEdges.overlapRect(screenEdges)) {
                    overlapping.push(screen);
                }
            }
        }

        return overlapping;
    }

    /**
     * Get an array of screens that touch edge to edge with this group
     */
    findTouchingScreens(group) {
        let overlapping = [];
        let groupEdges = edges(group, { pos: Vector.zero });

        for (let g = 0; g < this.alpine.groups.length; g++) {
            if (this.alpine.groups[g].id == group.id) {
                continue;
            }

            for (let s = 0; s < this.alpine.groups[g].screens.length; s++) {
                let screen = this.alpine.groups[g].screens[s];
                let screenEdges = edges(screen, this.alpine.groups[g]);

                if (groupEdges.touchesRect(screenEdges)) {
                    overlapping.push({ 
                        screen, 
                        edge: this._findClosestEdge(group, screen, this.alpine.groups[g])
                    });
                }
            }
        }

        return overlapping;
    }

    /**
     * Calculate a new set of transition zones based on the current layout
     */
    calculateTransitionZones() {
        let zones = {};

        for (let g = 0; g < this.alpine.groups.length; g++) {
            let group = this.alpine.groups[g];
            let overlaps = this.findOverlappingScreens(group);
            if (!overlaps.length) {
                continue;
            }

            for (let i = 0; i < overlaps.length; i++) {
                let screenGroup = this.alpine.findNeighbour(overlaps[i].groupId);
                let reverseOverlaps = this.findOverlappingScreens(screenGroup);

                if (!reverseOverlaps.length) {
                    continue;
                }

                if (reverseOverlaps.filter(overlap => overlap.groupId == group.id).length) {
                    return;
                }
            }
        }

        for (let g = 0; g < this.alpine.groups.length; g++) {
            let group = this.alpine.groups[g];
            let touches = this.findTouchingScreens(group);
            console.log(this.findTouchingScreens(group));
            zones[group.id] = [];

            for (let i = 0; i < touches.length; i++) {
                let { screen, edge } = touches[i];

                let direction = Directions.fromRect(edge[0]);
                zones[group.id].push({
                    id: screen.groupId,
                    bounds:  bounds(direction, group, screen, this.alpine.findNeighbour(screen.groupId)),
                    direction 
                })
            }
        }

        return zones;
    }

    _handleDragEnd(e) {
        e.preventDefault();
        document.onmousemove = null;

        let group = this.target;
        this.centerCanvas();
        console.log(this.findOverlappingScreens(group));
        console.log(this.findTouchingScreens(group));
        console.log(this.calculateTransitionZones());

        this.startPos = null;
        this.currPos = null;
        this.target = null;
    }

    _handleDragMove(e) {
        e.preventDefault()

        let pos = new Vector(e.clientX, e.clientY);
        let delta = this.currPos.subtract(pos);

        this.currPos = pos;

        this.target.pos = this.target.pos.subtract(delta);

        // TODO: this needs to work off overlapping edges rather than just using closestScreen.closestEdge
        this._snap(this.target);
    }

    _snap(group) {
        let screen = this._findClosestScreen(group);
        if (!screen) {
            return;
        }

        let screenGroup = this.alpine.findNeighbour(screen.groupId);
        if (!screenGroup) {
            return;
        }

        let edges = this._findClosestEdge(group, screen, screenGroup);
        if (!edges) {
            return;
        }

        let newPos = new Vector(group.pos.x, group.pos.y);

        if (edges[0] == "y") {
            newPos.y = screenGroup.pos.y + screen.pos.y + screen.height;
        } else if (edges[0] == "w") {
            newPos.x = screenGroup.pos.x + screen.pos.x - group.width;
        } else if (edges[0] == "z") {
            newPos.y = screenGroup.pos.y + screen.pos.y - group.height;
        } else {
            newPos.x = screenGroup.pos.x + screen.pos.x + screen.width;
        }

        if (newPos.distance(group.pos) <= SNAP_THRESHOLD) {
            group.pos = newPos;
            group.time = +new Date();
        }
    }

    _findClosestScreen(group) {
        let closest = null;
        let minDistance = 1 << 16;
        let groupCorners = corners(group, { pos: Vector.zero });

        for (let i = 0; i < this.alpine.groups.length; i++) {
            let target  = this.alpine.groups[i];
            if (target.id == group.id) {
                continue;
            }

            for (let s = 0; s < target.screens.length; s++) {
                let screenCorners = corners(target.screens[s], target);

                for (let x = 0; x < 4; x++)
                for (let y = 0; y < 4; y++) {
                    let distance = groupCorners[x].distance(screenCorners[y]);

                    if (distance < minDistance) {
                        minDistance = distance;
                        closest = target.screens[s];
                    }
                }
            }
        }

        return closest;
    }

    _findClosestEdge(group, screen, screenGroup) {
        let minDistance = 1 << 16;
        let closest = null;

        let groupEdges = edges(group, { pos: Vector.zero });
        let screenEdges = edges(screen, screenGroup);

        for (let i = 0; i < edgeChecks.length; i++) {
            let [ groupEdge, screenEdge ] = edgeChecks[i];

            let distance = Math.abs(groupEdges[groupEdge] - screenEdges[screenEdge]);
            if (distance < minDistance) {
                minDistance = distance;
                closest = edgeChecks[i];
            }
        }

        return closest;
    }
}

const corners = (screen, group) => {
    return [
        screen.pos.add(group.pos),
        new Vector(screen.pos.x + screen.width, screen.pos.y).add(group.pos),
        screen.pos.add(new Vector(screen.width, screen.height)).add(group.pos),
        new Vector(screen.pos.x, screen.pos.y + screen.height).add(group.pos)
    ];
};

const edges = (screen, group) => {
    return Vector4.fromRect(
        group.pos.x + screen.pos.x,
        group.pos.y + screen.pos.y,
        screen.width,
        screen.height
    );
}

const bounds = (direction, group, screen, screenGroup) => {
    switch (direction) {
        case Directions.UP:
            return [
                new Vector(
                    Math.max(group.pos.x, screenGroup.pos.x + screen.pos.x),
                    0
                ),
                new Vector(
                    Math.min(group.pos.x + group.width,  screenGroup.pos.x + screen.pos.x + screen.width),
                    0
                )
            ];

        case Directions.RIGHT:
            return [
                new Vector(
                    group.width,
                    Math.max(0, screenGroup.pos.y - group.pos.y + screen.pos.y)
                ),
                new Vector(
                    group.width,
                    Math.min(group.height, screenGroup.pos.y - group.pos.y + screen.pos.y + screen.height)
                )
            ];

        case Directions.DOWN:
            return [
                new Vector(
                    Math.max(group.pos.x, screenGroup.pos.x + screen.pos.x),
                    group.height
                ),
                new Vector(
                    Math.min(group.pos.x + group.width,  screenGroup.pos.x + screen.pos.x + screen.width),
                    group.height
                )
            ];

        case Directions.LEFT:
            return [
                new Vector(
                    0,
                    Math.max(0, screenGroup.pos.y - group.pos.y + screen.pos.y)
                ),
                new Vector(
                    0,
                    Math.min(group.height, screenGroup.pos.y - group.pos.y + screen.pos.y + screen.height)
                )
            ];
    }
};

const edgeChecks = [
    ["x", "w"],
    ["w", "x"],
    ["y", "z"],
    ["z", "y"],
];

export default ScreenMover;
