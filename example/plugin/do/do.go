package do

func serviceCallOne() (string, error) {
	return "a", nil
}

func serviceCallTwo() (int, error) {
	return 1, nil
}

func serviceCalls() (string, int, error) {
	return deriveDo(serviceCallOne, serviceCallTwo)
}
