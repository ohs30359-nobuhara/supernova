package template

type Result struct {
	Status Status
	Body   string
	Err    error
}

type Status int

const (
	SUCCESS Status = iota
	WARNING
	DANGER
)

func NewResultError(body string, status Status, err error) Result {
	return Result{
		Body:   body,
		Status: status,
		Err:    err,
	}
}

func NewResultSuccess(body string) Result {
	return Result{
		Body:   body,
		Status: SUCCESS,
		Err:    nil,
	}
}

/*
こんかなんじのエラーメッぜー

■ error description
```

```


*/
