package core

type Instance interface {
    ID() (ID string)
    Start()
    Stop()
    Pause()
    Kill()
}
