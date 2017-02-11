package validation

func VerifyTask(name string) {
	if err := verifyTask(name); err != nil {
		panic(err)
	}
}

func VerifyProject(name string) {
	if err := verifyProject(name); err != nil {
		panic(err)
	}
}

func VerifyTag(name string) {
	if err := verifyTag(name); err != nil {
		panic(err)
	}
}

func VerifyGoal(name string) {
	if err := verifyGoal(name); err != nil {
		panic(err)
	}
}

func VerifyHabit(name string) {
	if err := verifyHabit(name); err != nil {
		panic(err)
	}
}
