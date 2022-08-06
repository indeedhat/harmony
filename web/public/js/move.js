import Vector from '/js/vector.js';

class ScreenMover {
    constructor(harmony) {
        this.harmony = harmony;
        this.screens = document.querySelectorAll(".screen");
        this.groups = document.querySelectorAll(".screen-group");

        this.target = null;
        this.startPos = null;
        this.currPos = null;

        document.onmouseup = this._handleDragEnd.bind(this);
    }

    handleDragStart(e, group) {
        console.log("drag start");

        this.target = group;
        this.startPos = new Vector(e.clientX, e.clientY);
        this.currPos = new Vector(e.clientX, e.clientY);

        document.onmousemove = this._handleDragMove.bind(this);
    }

    _handleDragEnd() {
        document.onmousemove = null;
        console.log("drag end");

        // TODO: edge snapping

        this.startPos = null;
        this.currPos = null;
    }

    _handleDragMove(e) {
        e.preventDefault()
        console.log("drag");

        let pos = new Vector(e.clientX, e.clientY);
        let delta = this.currPos.subtract(pos);

        this.currPos = pos;

        this.target.pos = this.target.pos.subtract(delta);
    }
}

export default ScreenMover;
