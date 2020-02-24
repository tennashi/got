package got

func get(ioStream *ioStream, pkgName, cmdName string, isUpdate bool) error {
	pkgName = fullPackageName(pkgName, cmdName)
	dataDir, err := ensureDataDir()
	if err != nil {
		return err
	}
	cmdCtx := newCommandContext(ioStream, dataDir)

	if err := ensureGoModInit(cmdCtx, "got_local"); err != nil {
		return err
	}
	if err := goGet(cmdCtx, pkgName, isUpdate); err != nil {
		return err
	}
	if err := appendImport(dataDir, pkgName); err != nil {
		return err
	}
	if err := goModTidy(cmdCtx); err != nil {
		return err
	}
	return nil
}

func getAll(ioStream *ioStream, isUpdate bool) error {
	dataDir, err := ensureDataDir()
	if err != nil {
		return err
	}

	imports, err := getImports(dataDir)
	if err != nil {
		return err
	}

	for _, path := range imports {
		if err := get(ioStream, path, "", isUpdate); err != nil {
			return err
		}
	}
	return nil
}
