package utils

type Word struct {
    Raw string
    Type TokenType
    Line int
    Start int
    End int
}
