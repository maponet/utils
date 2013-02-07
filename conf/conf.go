/*
Copyright 2013 Petru Ciobanu, Francesco Paglia, Lorenzo Pierfederici

This file is part of Mapo.

Mapo is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 2 of the License, or
(at your option) any later version.

Mapo is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Mapo.  If not, see <http://www.gnu.org/licenses/>.
*/

/*
Package admin implements the API for Mapo's administration components.
*/
package conf

import (
	"goconf/conf"
)

/*
GlobalConfiguration, il oggetto globale per l'accesso ai dati contenuti nel
file di configurazione.
*/
var GlobalConfiguration *conf.ConfigFile

/*
ReadConfiguration, attiva il GlobalConfiguration.
*/
func ParseConfigFile(filepath string) error {

	c, err := conf.ReadConfigFile(filepath)
	if err == nil {
		GlobalConfiguration = c
	}

	return err
}
