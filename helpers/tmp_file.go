package helpers

import (
	"github.com/rs/zerolog/log"
	"os"
)

func GetTmpFile(data []byte) *os.File {
	tmpFile, err := os.CreateTemp(os.TempDir(), "brambu-telegram-bot")
	if err != nil {
		log.Error().Err(err).Msg("could not create tmpFile")
		return tmpFile
	}
	if _, err = tmpFile.Write(data); err != nil {
		log.Error().Err(err).Msg("could not write to tmpFile")
		return tmpFile
	}
	if _, err = tmpFile.Seek(0, 0); err != nil {
		log.Error().Err(err).Msg("could not rewind tmpFile")
		return tmpFile
	}
	return tmpFile
}

func CleanupTmpFile(tmpFile *os.File) {
	if tmpFile == nil {
		return
	}
	if err := os.RemoveAll(tmpFile.Name()); err != nil {
		log.Error().Err(err).Msg("failed to remove temporary file")
	}
}
