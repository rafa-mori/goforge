package article

// This file/package allows the article module to be used as a library.
// It defines the ArticleMod interface which can be implemented by any module
// that wants to be part of the article ecosystem.

type ArticleMod interface {
	// Active returns true if the module is active.
	Active() bool
	// Module returns the name of the module.
	Module() string
	// Execute runs the module's command.
	Execute() error
}
