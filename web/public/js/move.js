import Vector from '/js/vector.js';

class ScreenMover {
    constructor(alpine) {
        this.alpine = alpine;

        this.screens = document.querySelectorAll(".screen");
        this.groups = document.querySelectorAll(".screen-group");

        this.target = null;
        this.startPos = null;
        this.currPos = null;

        document.onmouseup = this._handleDragEnd.bind(this);
    }

    handleDragStart(e, group) {
        this.target = group;
        this.startPos = new Vector(e.clientX, e.clientY);
        this.currPos = new Vector(e.clientX, e.clientY);

        document.onmousemove = this._handleDragMove.bind(this);
    }

    _handleDragEnd() {
        // TODO: this doesnt get refleted until the next proxy change for some reason
        this._snap(this.target);

        document.onmousemove = null;

        // TODO: reposition whole canvas

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
    }

    _snap(group) {
        let screen = this._findClosestScreen(group);
        if (!screen) {
            return;
        }

        let edges = this._findClosestEdge(group, screen);
        if (!edges) {
            return;
        }

        let screenGroup = this.alpine.findNeighbour(screen.groupId);
        if (!screenGroup) {
            return;
        }

        let newPos = new Vector(group.pos.x, group.pos.y);

        if (edges[0] == "top") {
            newPos.y = screenGroup.pos.y + screen.pos.y + group.height;
        } else if (edges[0] == "right") {
            newPos.x = screenGroup.pos.x + screen.pos.x - group.width;
        } else if (edges[0] == "bottom") {
            newPos.y = screenGroup.pos.y + screen.pos.y - group.height;
        } else {
            newPos.x = screenGroup.pos.x + screen.pos.x + screen.width;
        }

        console.log({
            orig: group.pos.add(Vector.zero),
            new: newPos
        })

        group.pos = newPos;
        group.time = +new Date();
    }

    _findClosestScreen(group) {
        let closest = null;
        let minDistance = 1 << 16;
        let groupCorners = corners(group);

        for (let i = 0; i < this.alpine.groups.length; i++) {
            let target  = this.alpine.groups[i];
            if (target.id == group.id) {
                continue;
            }

            for (let s = 0; s < target.screens.length; s++) {
                let screenCorners = corners(target.screens[s]);

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

    _findClosestEdge(group, screen) {
        console.log({ screen, group })
        let minDistance = 1 << 16;
        let closest = null;

        let groupEdges = edgeCenters(group);
        let screenEdges = edgeCenters(screen);

        for (let i = 0; i < edgeChecks.length; i++) {
            let [ groupEdge, screenEdge ] = edgeChecks[i];

            let distance = groupEdges[groupEdge].distance(screenEdges[screenEdge]);
            console.log(groupEdges[groupEdge], screenEdges[screenEdge], distance, minDistance)
            if (distance < minDistance) {
                minDistance = distance;
                closest = edgeChecks[i];
            }
        }

        return closest;
    }
}


const corners = screen => {
    return [
        screen.pos,
        new Vector(screen.pos.x + screen.width, screen.pos.y),
        screen.pos.add(new Vector(screen.width, screen.height)),
        new Vector(screen.pos.x, screen.pos.y + screen.height)
    ];
};

const edgeCenters = screen => {
    return {
        top: new Vector(screen.pos.x + screen.width/2, screen.pos.y),
        right: new Vector(screen.pos.x + screen.width, screen.pos.y + screen.height/2),
        bottom: new Vector(screen.pos.x + screen.width/2, screen.pos.y + screen.height),
        left: new Vector(screen.pos.x, screen.pos.y + screen.height/2)
    }
}

const edgeChecks = [
    ["top", "bottom"],
    ["bottom", "top"],
    ["left", "right"],
    ["right", "left"],
];

export default ScreenMover;
