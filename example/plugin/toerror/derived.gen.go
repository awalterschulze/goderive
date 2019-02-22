// Code generated by goderive DO NOT EDIT.

package toerror

// deriveToError transforms the given function's last bool type into an error type. The transformed function returns the given error when the result of the given function is false, otherwise it returns nil.
func deriveToError(err error, f func(vers string) (major int, minor int, ok bool)) func(vers string) (int, int, error) {
	return func(vers string) (int, int, error) {
		out0, out1, success := f(vers)
		if success {
			return out0, out1, nil
		}
		return out0, out1, err
	}
}

// deriveCompose composes functions f0 and f1 into one function, that takes the parameters from f0 and returns the results from f1.
func deriveCompose(f0 func(string) (int, int, error), f1 func(int, int) (int, error)) func(string) (int, error) {
	return func(v_0_0 string) (int, error) {
		v_1_0, v_1_1, err0 := f0(v_0_0)
		if err0 != nil {
			return 0, err0
		}
		v_2_0, err1 := f1(v_1_0, v_1_1)
		if err1 != nil {
			return 0, err1
		}
		return v_2_0, nil
	}
}
