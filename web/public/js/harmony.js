const Harmony = groups => {
    console.log("called")
    Alpine.data("groups", () => ({
        init() {
            console.log(this.groups);
        },

        groups: format(groups),
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

    return {
        pos: new Vector(),
        width: data.Width,
        height: data.Height,
        screens
    };
};

const screen = data => ({
    pos: new Vector(data.Position.X, data.Position.Y),
    width: data.Width,
    height: data.Height
});

class Vector {
    constructor(x = 0, y = 0) {
        this.x = x;
        this.y = y;
    }
}

export default Harmony;
