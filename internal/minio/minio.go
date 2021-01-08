/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 08/01/2021
*/

package minio

import (
	"fmt"
	"time"
)

// PutFile get a file by filename and upload it in minio server
func PutFile(, source, target string) error {
		start := time.Now()
		defer addons.MethodMetricsTime("minio", "putfile", start)
		log.WithFields(loguuid).Debug("[start]")
		defer log.WithFields(loguuid).Debug("[end]")

		if s.Env != "test" {
		mc := s.getMinioClient(loguuid)
		if mc == nil {
		log.WithFields(loguuid).Error("generate minio client failed")
		return fmt.Errorf("backup failed")
	}
		log.WithFields(loguuid).Debugf("put %s to artifacts (%s)", source, target)
		if _, err := mc.FPutObject(s.Minio.Bucket, target, source, minio.PutObjectOptions{}); err != nil {
		log.WithFields(loguuid).Error(err)
		return fmt.Errorf("backup failed")
	}
		log.WithFields(loguuid).Debugf("File %s uploaded", source)
	} else {
		return mock.PutToMinio(target)
	}
		return nil
	}
}
