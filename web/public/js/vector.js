class Vector {
    static zero = new Vector();

    static fromGo(goVector) {
        return new Vector(goVector.X, goVector.Y);
    }

    static topLeft(a, b) {
        if (a.distance(Vector.zero) < b.distance(Vector.zero)) {
            return a;
        }

        return b;
    }

    static overlapRect(a, b) {
        return a.left.x <= b.right.x
            &&  b.left.x <= a.right.x
            && a.top.y <= b.top.y
            && b.top.y <= a.top.y;
    }

    constructor(x = 0, y = 0) {
        this.x = x;
        this.y = y;
    }

    distance(pos) {
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

class Vector4 {
    constructor(x, y, w, z) {
        this.x = x;
        this.y = y;
        this.w = w;
        this.z = z;
    }

    static fromRect(x, y, width, height) {
        return new Vector4(x, y, x + width, y + height)
    }

    overlapRect(r2) {
        // no horizontal overlap
        if (this.x >= r2.w || r2.x >= this.w) {
            return false;
        }

        // no vertical overlap
        console.log(this.y, r2.z, r2.y, this.z)
        if (this.y >= r2.z || r2.y >= this.z) {
            return false;
        }

        return true;
    }

    touchesRect(r2) {
        // no horizontal overlap
        if (this.x > r2.w || r2.x > this.w) {
            return false;
        }

        // no vertical overlap
        if (this.y > r2.z || r2.y > this.z) {
            return false;
        }

        return true;
    }
}

export default Vector;
export {
    Vector4
};
