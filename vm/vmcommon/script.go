package vmcommon

func BuildP2PKHScript(publicKeyHash []byte) []byte {
	var script []byte
	script = append(script, byte(OpDup))
	script = append(script, byte(OpHash256))
	script = append(script, byte(OpPushData32))
	script = append(script, publicKeyHash...)
	script = append(script, byte(OpEqualVerify))
	script = append(script, byte(OpCheckSig))
	return script
}
