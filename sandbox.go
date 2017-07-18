package gisp

// Box which contains the raw userspace values
type Box map[string]interface{}

// Sandbox sandbox
// It implements prototype design pattern
type Sandbox struct {
	dict   Box
	parent *Sandbox
}

// New create a new sandbox
func New(dict Box) *Sandbox {
	return &Sandbox{
		dict: dict,
	}
}

// Create create a new sandbox which dirives from current sandbox
func (sandbox *Sandbox) Create() *Sandbox {
	return &Sandbox{
		dict:   Box{},
		parent: sandbox,
	}
}

// Get get property from the prototype chain
func (sandbox *Sandbox) Get(name string) (interface{}, bool) {
	for sandbox != nil {
		val, has := sandbox.dict[name]

		if has {
			return val, true
		}

		sandbox = sandbox.parent
	}

	return nil, false
}

// Names recursively get all names from current node to all its ancestors
func (sandbox *Sandbox) Names() []string {
	names := []string{}

	for sandbox != nil {
		for k := range sandbox.dict {
			names = append(names, k)
		}

		sandbox = sandbox.parent
	}

	return names
}

// Box return flat dict
func (sandbox *Sandbox) Box() Box {
	box := Box{}

	for sandbox != nil {
		for k := range sandbox.dict {
			box[k], _ = sandbox.Get(k)
		}

		sandbox = sandbox.parent
	}

	return box
}

// Set set property
func (sandbox *Sandbox) Set(name string, val interface{}) {
	sandbox.dict[name] = val
}

// Reset set property
// Update the nearest one to root, if nothing found a new property
// will be created on current closure
func (sandbox *Sandbox) Reset(name string, val interface{}) {
	curr := sandbox

	for sandbox != nil {
		_, has := sandbox.dict[name]

		if has {
			sandbox.dict[name] = val
			return
		}

		sandbox = sandbox.parent
	}

	curr.dict[name] = val
}
