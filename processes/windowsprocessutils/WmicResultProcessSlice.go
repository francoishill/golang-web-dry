package windowsprocessutils

type WmicResultProcessSlice []*WmicResultProcess

func (w WmicResultProcessSlice) FindFirstWithParentPid(parentPidToFind int) *WmicResultProcess {
	for _, rp := range w {
		if rp.ParentProcessId == parentPidToFind {
			return rp
		}
	}
	return nil
}
