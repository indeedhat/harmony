import Vector from '/js/vector.js';
import ScreenMover from '/js/move.js';

const Harmony = groups => {
    Alpine.data("groups", () => ({
        init() {
            this.reposition();
            this.mover = new ScreenMover(this);
        },

        handleDragStart(e, group) {
            console.log({ e, group })
            this.mover.handleDragStart(e, group);
        },

        reposition() {
            let done = [this.groups[0].id];

            for (let i = 0; i < this.groups.length; i++) {
                let group = this.groups[i];
                console.log({ group })

                for (let j = 0; j < group.transitions.length; j++) {
                    let trans = group.transitions[j];
                    if (~done.indexOf(trans.id)) {
                        continue;
                    }

                    let neighbour = this.findMatchingTransition(group.id, trans.id);
                        console.log({ neighbour })
                    if (!neighbour) {
                        continue;
                    }

                    let newPos = new Vector(
                        group.pos.x + trans.pos.x - neighbour.pos.x,
                        group.pos.y + trans.pos.y - neighbour.pos.y
                    );

                    let target = this.findNeighbour(trans.id);
                    if (target) {
                        target.pos = newPos;
                        done.push(target.id);
                    }
                }
            }
        },

        findMatchingTransition(self, neighbour) {
            for (let i = 0; i < this.groups.length; i++) {
                let group = this.groups[i];
                if (group.id != neighbour) {
                    continue;
                }

                for (let j = 0; j < group.transitions.length; j++) {
                    if (group.transitions[j].id == self) {
                        return group.transitions[j];
                    }
                }
            }

            return null;
        },

        findNeighbour(id) {
            for (let i = 0; i < this.groups.length; i++) {
                if (id == this.groups[i].id) {
                    return this.groups[i];
                }
            }

            return null;
        },

        groups: format(groups),
        mover: null
    }));

};

const format = data => {
    let groups = [];

    for (let i = 0; i < data.length; i++) {
        groups.push(group(data[i]));
    }

    return groups;
};

const group = data => {
    let screens = [];
    for (let i = 0; i < data.Displays.length; i++) {
        screens.push(screen(data.Displays[i]));
    }

    let transitions = [];
    for (let i = 0; i < (data.Transitions || []).length; i++) {
        transitions.push(transitionZone(data.Transitions[i]));
    }

    return {
        id: data.UUID,
        name: data.Hostname,
        pos: new Vector(),
        width: data.Width,
        height: data.Height,
        transitions,
        screens
    };
};

const screen = data => ({
    pos: new Vector(data.Position.X, data.Position.Y),
    width: data.Width,
    height: data.Height
});

const transitionZone = data => ({
    id: data.UUID,
    pos: Vector.topLeft(
        Vector.fromGo(data.Bounds[0]), 
        Vector.fromGo(data.Bounds[1])
    )
})


export default Harmony;
