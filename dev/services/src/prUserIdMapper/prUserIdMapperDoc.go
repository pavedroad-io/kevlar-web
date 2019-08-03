// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Microservice for mapping multiple 3rd party id into internal ID
//
// Given a key, login, or email address return the PR UUID
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: api.pavedroad.io
//     BasePath: /api/v1/namespace/pavedroad.io/prUserIdMappers
//     Version: 0.0.1
//     License: Apache 2
//     Contact: John Scharber<john@pavedroad.io> http://john.pavedroad.io
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main
