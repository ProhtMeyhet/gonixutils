package signal

import(
	"errors"
	"strings"
	"syscall"

	"github.com/ProhtMeyhet/libgosimpleton/simpleton"
)

func IsTerminal(signal syscall.Signal) bool {
	switch signal {
	case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGINT:
		return true
	}

	return false
}

func SignalToString(signal syscall.Signal) (stringSignal string) {
	stringSignal = "SIG"

	switch signal {
		case syscall.SIGHUP:
			stringSignal += "HUP"
		case syscall.SIGINT:
			stringSignal += "INT"
		case syscall.SIGQUIT:
			stringSignal += "QUIT"
		case syscall.SIGILL:
			stringSignal += "ILL"
		case syscall.SIGTRAP:
			stringSignal += "TRAP"
		case syscall.SIGABRT:
			stringSignal += "ABRT"
		case syscall.SIGBUS:
			stringSignal += "BUS"
		case syscall.SIGFPE:
			stringSignal += "FPE"
		case syscall.SIGKILL:
			stringSignal += "KILL"
		case syscall.SIGUSR1:
			stringSignal += "USR1"
		case syscall.SIGSEGV:
			stringSignal += "SEGV"
		case syscall.SIGUSR2:
			stringSignal += "USR2"
		case syscall.SIGPIPE:
			stringSignal += "PIPE"
		case syscall.SIGALRM:
			stringSignal += "ALRM"
		case syscall.SIGTERM:
			stringSignal += "TERM"
		case syscall.SIGSTKFLT:
			stringSignal += "STKFLT"
		case syscall.SIGCHLD:
			stringSignal += "CHLD"
		case syscall.SIGCONT:
			stringSignal += "CONT"
		case syscall.SIGSTOP:
			stringSignal += "STOP"
		case syscall.SIGTSTP:
			stringSignal += "TSTP"
		case syscall.SIGTTIN:
			stringSignal += "TTIN"
		case syscall.SIGTTOU:
			stringSignal += "TTOU"
		case syscall.SIGURG:
			stringSignal += "URG"
		case syscall.SIGXCPU:
			stringSignal += "XCPU"
		case syscall.SIGXFSZ:
			stringSignal += "XFSZ"
		case syscall.SIGVTALRM:
			stringSignal += "VTALRM"
		case syscall.SIGPROF:
			stringSignal += "PROF"
		case syscall.SIGWINCH:
			stringSignal += "WINCH"
		case syscall.SIGIO:
			stringSignal += "IO"
		case syscall.SIGPWR:
			stringSignal += "PWR"
		case syscall.SIGSYS:
			stringSignal += "SYS"
	}

	return
}

func StringToSignal(signal string) (syscall.Signal, error) {
	signal = strings.ToUpper(signal)
	signal = simpleton.RemoveFirstCharsIfEqual(signal, "SIG", "SIGNAL")

	switch strings.ToUpper(signal) {
		case "HUP":
			return syscall.SIGHUP, nil
		case "INT", "INTERRUPT":
			return syscall.SIGINT, nil
		case "Q", "QUIT":
                        return syscall.SIGQUIT, nil
                case "ILL", "ILLEGAL":
                        return syscall.SIGILL, nil
                case "TRAP":
                        return syscall.SIGTRAP, nil
                case "ABRT", "ABORT":
                        return syscall.SIGABRT, nil
                case "BUS":
                        return syscall.SIGBUS, nil
                case "FPE":
                        return syscall.SIGFPE, nil
                case "KILL":
                        return syscall.SIGKILL, nil
                case "USR1", "1", "USER1":
                        return syscall.SIGUSR1, nil
                case "SEGV", "SEGFAULT":
                        return syscall.SIGSEGV, nil
                case "USR2", "2", "USER2":
                        return syscall.SIGUSR2, nil
                case "PIPE":
                        return syscall.SIGPIPE, nil
                case "ALRM", "ALARM":
                        return syscall.SIGALRM, nil
                case "TERM", "TERMINATE":
                        return syscall.SIGTERM, nil
                case "STKFLT":		// coprocessor fault
                        return syscall.SIGSTKFLT, nil
                case "CHLD", "CHILD":
                        return syscall.SIGCHLD, nil
                case "CONT", "CONTINUE":
                        return syscall.SIGCONT, nil
                case "STOP":
                        return syscall.SIGSTOP, nil
                case "STP":
                        return syscall.SIGTSTP, nil
                case "TTIN":
                        return syscall.SIGTTIN, nil
                case "TTOU":
                        return syscall.SIGTTOU, nil
                case "URG", "URGENT":
                        return syscall.SIGURG, nil
                case "XCPU":
                        return syscall.SIGXCPU, nil
                case "XFSZ":
                        return syscall.SIGXFSZ, nil
                case "VTALRM":
                        return syscall.SIGVTALRM, nil
                case "PROF":
                        return syscall.SIGPROF, nil
                case "WINCH":		// window changed
                        return syscall.SIGWINCH, nil
                case "IO":
                        return syscall.SIGIO, nil
                case "PWR":
                        return syscall.SIGPWR, nil
                case "SYS":
                        return syscall.SIGSYS, nil
	}

	return syscall.SIGTERM, errors.New("invalid signal!")
}

func IntSignalToString(signal int) string {
	switch signal {
		case 1:
			return "SIGHUP"
		case 2:
			return "SIGINT"
		case 3:
			return "SIGQUIT"
		case 4:
			return "SIGILL"
		case 5:
			return "SIGTRAP"
		case 6:
			return "SIGABRT"
		case 7:
			return "SIGBUS"
		case 8:
			return "SIGFPE"
		case 9:
			return "SIGKILL"
		case 10:
			return "SIGUSR1"
		case 11:
			return "SIGSEGV"
		case 12:
			return "SIGUSR2"
		case 13:
			return "SIGPIPE"
		case 14:
			return "SIGALRM"
		case 15:
			return "SIGTERM"
		case 16:
			return "SIGSTKFLT"
		case 17:
			return "SIGCHLD"
		case 18:
			return "SIGCONT"
		case 19:
			return "SIGSTOP"
		case 20:
			return "SIGTSTP"
		case 21:
			return "SIGTTIN"
		case 22:
			return "SIGTTOU"
		case 23:
			return "SIGURG"
		case 24:
			return "SIGXCPU"
		case 25:
			return "SIGXFSZ"
		case 26:
			return "SIGVTALRM"
		case 27:
			return "SIGPROF"
		case 28:
			return "SIGWINCH"
		case 29:
			return "SIGIO"
		case 30:
			return "SIGPWR"
		case 31:
			return "SIGSYS"
		case 34:
			return "SIGRTMIN"
		case 35:
			return "SIGRTMIN"
		default:
			return INVALID
	}
}

func IntToSignal(signal int) (syscall.Signal, error) {
	switch signal {
		case 1:
			return syscall.SIGHUP, nil
		case 2:
			return syscall.SIGINT, nil
		case 3:
			return syscall.SIGQUIT, nil
		case 4:
			return syscall.SIGILL, nil
		case 5:
			return syscall.SIGTRAP, nil
		case 6:
			return syscall.SIGABRT, nil
		case 7:
			return syscall.SIGBUS, nil
		case 8:
			return syscall.SIGFPE, nil
		case 9:
			return syscall.SIGKILL, nil
		case 10:
			return syscall.SIGUSR1, nil
		case 11:
			return syscall.SIGSEGV, nil
		case 12:
			return syscall.SIGUSR2, nil
		case 13:
			return syscall.SIGPIPE, nil
		case 14:
			return syscall.SIGALRM, nil
		case 15:
			return syscall.SIGTERM, nil
		case 16:
			return syscall.SIGSTKFLT, nil
		case 17:
			return syscall.SIGCHLD, nil
		case 18:
			return syscall.SIGCONT, nil
		case 19:
			return syscall.SIGSTOP, nil
		case 20:
			return syscall.SIGTSTP, nil
		case 21:
			return syscall.SIGTTIN, nil
		case 22:
			return syscall.SIGTTOU, nil
		case 23:
			return syscall.SIGURG, nil
		case 24:
			return syscall.SIGXCPU, nil
		case 25:
			return syscall.SIGXFSZ, nil
		case 26:
			return syscall.SIGVTALRM, nil
		case 27:
			return syscall.SIGPROF, nil
		case 28:
			return syscall.SIGWINCH, nil
		case 29:
			return syscall.SIGIO, nil
		case 30:
			return syscall.SIGPWR, nil
		case 31:
			return syscall.SIGSYS, nil
	}

	return syscall.SIGTERM, errors.New("invalid signal!")
}
