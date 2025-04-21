// AGITA - Atlassian-Github Issue Transfer Application

package main

func main() {
    GetArgs()
    switch Mode {
        case ModeTransfer:  TransferAll(Args...)
        case ModeExport:    ExportAll(Args...)
        case ModeClear:     ClearAll(Args...)
        case ModeTrial:     TrialAll(Args...)
        default:            panic("main action undefined")
    }
}
