#!/bin/bash -eu
#
# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/bash

final_zip=sqlcommenter-mono.zip
original_branch=$(git rev-parse --abbrev-ref HEAD)
gen_branch=expand-$(date | tr " " "_" | tr ":" "_")
git checkout -b $gen_branch || (echo "Failed to checkout $gen_branch" && exit)

rm -rf $final_zip && zip -r $final_zip $(ls) || (git checkout $original_branch && git branch -D $gen_branch && exit)
git checkout $original_branch && git branch -D $gen_branch
