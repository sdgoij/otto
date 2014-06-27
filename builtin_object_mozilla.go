package otto

func builtinObject___defineGetter__(call FunctionCall) Value {
	if fn := call.Argument(1); fn.isCallable() {
		this, key := call.thisObject(), call.Argument(0).String()
		getset := _propertyGetSet{fn._object(), nil}
		if property := this.getProperty(key); nil != property {
			if current, test := property.value.(_propertyGetSet); test {
				getset[1] = current[1]
			}
		}
		descriptor := _property{getset, _propertyMode(0011)}
		this.defineOwnProperty(key, descriptor, false)
		return UndefinedValue()
	}
	panic(call.runtime.panicTypeError())
}

func builtinObject___defineSetter__(call FunctionCall) Value {
	if fn := call.Argument(1); fn.isCallable() {
		this, key := call.thisObject(), call.Argument(0).String()
		getset := _propertyGetSet{nil, fn._object()}
		if property := this.getProperty(key); nil != property {
			if current, test := property.value.(_propertyGetSet); test {
				getset[0] = current[0]
			}
		}
		descriptor := _property{getset, _propertyMode(0011)}
		this.defineOwnProperty(key, descriptor, false)
		return UndefinedValue()
	}
	panic(call.runtime.panicTypeError())
}

func builtinObject___lookupGetter__(call FunctionCall) Value {
	if property := call.thisObject().getProperty(call.Argument(0).String()); nil != property {
		if getset, test := property.value.(_propertyGetSet); test {
			return toValue(getset[0])
		}
	}
	return UndefinedValue()
}

func builtinObject___lookupSetter__(call FunctionCall) Value {
	if property := call.thisObject().getProperty(call.Argument(0).String()); nil != property {
		if getset, test := property.value.(_propertyGetSet); test {
			return toValue(getset[1])
		}
	}
	return UndefinedValue()
}
