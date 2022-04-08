<?php

use Google\GoogleSqlCommenterLaravel\Utils;
use PHPUnit\Framework\TestCase;

final class UtilsTest extends TestCase
{
    public function testformatcommentswithkeys(): void
    {
        $this->assertEquals("/*key1='value1',key2='value2'*/",Utils::format_comments(array("key1"=>"value1", "key2"=>"value2")));
    }
    public function testformatcommentswithoutkeys(): void
    {
        $this->assertEquals("",Utils::format_comments(array()));
    }
}