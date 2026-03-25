// ExitBox - Multi-Agent Container Sandbox
// Copyright (C) 2026 Cloud Exit B.V.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package agents

import (
	"testing"
)

func TestRegistryGetAll(t *testing.T) {
	for _, name := range Names() {
		a := Get(name)
		if a == nil {
			t.Errorf("Get(%q) returned nil", name)
			continue
		}
		if a.Name() != name {
			t.Errorf("Get(%q).Name() = %q", name, a.Name())
		}
	}
}

func TestGetUnknown(t *testing.T) {
	if a := Get("nonexistent"); a != nil {
		t.Errorf("Get(\"nonexistent\") should return nil, got %v", a)
	}
}
