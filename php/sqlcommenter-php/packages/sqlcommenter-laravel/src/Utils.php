<?php

namespace Google\GoogleSqlCommenterLaravel;

class Utils
{
    public static function format_comments($comment)
    {
        if (empty($comment)) {
            return "";
        }
        $lastElement = array_key_last($comment);
        $sql_comment = "/*";
        foreach($comment as $key=>$value)
        {
            if ($key == $lastElement){
                $sql_comment .=$key ."=". "'".$value."'*/";
            }
            else{
                $sql_comment .=$key ."=". "'".$value."',";
            }
        }
        return $sql_comment;
    }

}