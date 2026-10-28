module coursemod

go 1.26.0

// Add the library your course teaches here (go get <lib>@<ver>), run go mod tidy to create go.sum,
// and keep this go.mod/go.sum in src/_root/. build.py copies them into every task and replaces only
// the 'module coursemod' line with the task's module name.
