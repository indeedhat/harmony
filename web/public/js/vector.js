class Vector {
    static zero = new Vector();

    static fromGo(goVector) {
        return new Vector(goVector.X, goVector.Y);
    }

    static topLeft(a, b) {
        console.log({a, ad : a.distanceFrom(Vector.zero), b, bd: b.distanceFrom(Vector.zero)});
        if (a.distanceFrom(Vector.zero) < b.distanceFrom(Vector.zero)) {
            return a;
        }

        return b;
}

    constructor(x = 0, y = 0) {
        this.x = x;
        this.y = y;
    }

    distanceFrom(pos) {
        return Math.sqrt(
            Math.abs(
                (this.x - pos.x) ** 2 
                + (this.y - pos.y) ** 2
            )
        );
    }

    subtract(pos) {
        return new Vector(
            this.x - pos.x,
            this.y - pos.y
        );
    }

    add(pos) {
        return new Vector(
            this.x + pos.x,
            this.y + pos.y
        );
    }
}

export default Vector;
