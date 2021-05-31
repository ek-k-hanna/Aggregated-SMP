/*
 * Copyright (C) 2019 ING BANK N.V.
 *
 * May 2021 - Modified by Hanna Ek.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
 

package bulletproofs

import (
    "hannaekthesis/p256"
    _"math/big"

)
var ORDER = p256.CURVE.N
var SEEDH = "BulletproofsDoesNotNeedTrustedSetupH"
var MAX_RANGE_END int64 = 256//4294967296 // 2**32
var MAX_RANGE_END_EXPONENT = 8//32      // 2**32
