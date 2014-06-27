package otto

import (
	"testing"
)

func TestObject___defineGetter__(t *testing.T) {
	tt(t, func() {
		test, _ := test()
		test(`raise: this.__defineGetter__("a", "x")`, "TypeError")
		test(`
        var object = {}
        object.__defineGetter__("a", function() { return true })
        object.a
    `, true)
		test(`
        var object = {x: 10, y: 20, z: false};
        object.__defineGetter__("xy", function() { return this.x + this.y });
        object.__defineGetter__("z", function() { return true });
        [object.xy, object.z];
    `, "30,true")
	})
}

func TestObject___lookupGetter__(t *testing.T) {
	tt(t, func() {
		test, _ := test()
		test(`
        var object = {};
        var x = object.__lookupGetter__("a") === undefined;
        object.__defineGetter__("a", function() { return true });
        [x, object.a, object.__lookupGetter__("a") !== undefined];
    `, "true,true,true")
	})
}

func TestObject___defineSetter__(t *testing.T) {
	tt(t, func() {
		test, _ := test()
		test(`raise: this.__defineSetter__("a", "x")`, "TypeError")
		test(`
        var value = false
        var object = {}
        object.__defineSetter__("a", function(v) { value = v })
        object.a = true
        value === true
    `, true)
		test(`
        var object = {x: 20, y: 10, z: 30};
        object.__defineSetter__("xyz", function(v) {
            if (Array.isArray(v)) {
                this.x = v[0]
                this.y = v[1]
                this.z = v[2]
                return
            }
            throw new TypeError()
        })
        object.xyz = [30, 20, 10];
        var caught = false;
        try {
            object.xyz = "abc"
        } catch (e) {
            caught = e instanceof TypeError
        }
        [object.x, object.y, object.z, caught];
    `, "30,20,10,true")
	})
}

func TestObject___lookupSetter__(t *testing.T) {
	tt(t, func() {
		test, _ := test()
		test(`
        var value = false;
        var object = {};
        var x = object.__lookupSetter__("a") === undefined;
        object.__defineSetter__("a", function(v) { value = v });
        object.a = true;
        [x, value, object.__lookupSetter__("a") !== undefined];
    `, "true,true,true")
	})
}
