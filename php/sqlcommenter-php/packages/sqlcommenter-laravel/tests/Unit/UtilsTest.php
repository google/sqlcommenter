<?php
/*
 * Copyright 2022 Google LLC

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http:*www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

use Google\GoogleSqlCommenterLaravel\Utils;
use PHPUnit\Framework\TestCase;

final class UtilsTest extends TestCase
{
    public function testFormatCommentsWithKeys(): void
    {
        $this->assertEquals("/*key1='value1',key2='value2'*/", Utils::formatComments(["key1" => "value1", "key2" => "value2"]));
    }

    public function testFormatCommentsWithoutKeys(): void
    {
        $this->assertEquals("", Utils::formatComments([]));
    }

    public function testFormatCommentsWithSpecialCharKeys(): void
    {
        $this->assertEquals("/*key1='value1%%40',key2='value2'*/", Utils::formatComments(["key1" => "value1@", "key2" => "value2"]));
    }

    public function testFormatCommentsWithPlaceholder(): void
    {
        $this->assertEquals("/*key1='value1%%3F',key2='value2'*/", Utils::formatComments(["key1" => "value1?", "key2" => "value2"]));
    }

    public function testFormatCommentsWithNamedPlaceholder(): void
    {
        $this->assertEquals("/*key1='%%3Anamed',key2='value2'*/", Utils::formatComments(["key1" => ":named", "key2" => "value2"]));
    }
}
