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

namespace Google\GoogleSqlCommenterLaravel;

class Utils
{
    public static function formatComments($comment)
    {
        if (empty($comment)) {
            return "";
        }
        $lastElement = array_key_last($comment);
        $sql_comment = "/*";
        foreach ($comment as $key => $value) {
            if ($key == $lastElement) {
                $sql_comment .= Utils::customUrlEncode($key) . "=" . "'" . Utils::customUrlEncode($value) . "'*/";
            } else {
                $sql_comment .= Utils::customUrlEncode($key) . "=" . "'" . Utils::customUrlEncode($value) . "',";
            }
        }
        return $sql_comment;
    }

    private static function customUrlEncode($input)
    {
        $encodedString = urlencode($input);

        # Since SQL uses '%' as a keyword, '%' is a by-product of url quoting
        # e.g. foo,bar --> foo%2Cbar
        # thus in our quoting, we need to escape it too to finally give
        #      foo,bar --> foo%%2Cbar

        return str_replace("%", "%%", $encodedString);
    }
}
