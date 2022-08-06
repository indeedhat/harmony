import Vector from '/js/vector.js';
import ScreenMover from '/js/move.js';

const Harmony = groups => {
    Alpine.data("groups", () => ({
        init() {
            this.reposition();
            this.mover = new ScreenMover(this);
        },

        handleDragStart(e, group) {
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
    return data.map(group => ({
        id: group.UUID,
        time: +new Date(), // thest to force update
        name: group.Hostname,
        pos: new Vector(),
        width: group.Width,
        height: group.Height,
        screens: group.Displays.map(screen => ({
            groupId: group.UUID,
            pos: new Vector(screen.Position.X, screen.Position.Y),
            width: screen.Width,
            height: screen.Height
        })),
        transitions: group.Transitions.map(transition => ({
            id: transition.UUID,
            pos: Vector.topLeft(
                Vector.fromGo(transition.Bounds[0]), 
                Vector.fromGo(transition.Bounds[1])
            )
        }))
    }))
};

export default Harmony;
