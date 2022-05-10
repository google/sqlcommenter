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
namespace Tests\Feature;

use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;
use Illuminate\Support\Facades\DB;
use Config;

class DBTest extends TestCase
{
    public function test_delete()
    {
        DB::enableQueryLog();
        $response = $this->get('api/db/delete');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~delete from users where name='johnny'/\*framework='laravel-\d*.\d*.\d*',controller='RawDBController',action='delete',route='%%2Fapi%%2Fdb%%2Fdelete',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/;~", end($myDebugVar)['query']);
    }
    public function test_insert()
    {
        DB::enableQueryLog();
        $response = $this->get('api/db/insert');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~insert into users \(`name`, `email`, `password`\) values \('john', 'contact_me@daa.com', 'Passsword3456'\)/\*framework='laravel-\d*.\d*.\d*',controller='RawDBController',action='insert',route='%%2Fapi%%2Fdb%%2Finsert',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }

    public function test_select_one()
    {
        DB::enableQueryLog();
        $response = $this->get('api/db/selectone');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~select 1/\*framework='laravel-\d*.\d*.\d*',controller='RawDBController',action='selectOne',route='%%2Fapi%%2Fdb%%2Fselectone',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }

    public function test_update()
    {
        DB::enableQueryLog();
        $response = $this->get('api/db/update');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~update users set name='johnny' where name='john'/\*framework='laravel-\d*.\d*.\d*',controller='RawDBController',action='update',route='%%2Fapi%%2Fdb%%2Fupdate',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }

    public function test_select()
    {
        DB::enableQueryLog();
        $response = $this->get('api/db/select');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~select 1/\*framework='laravel-\d*.\d*.\d*',controller='RawDBController',action='select',route='%%2Fapi%%2Fdb%%2Fselect',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }

    public function test_select_with_route_disabled()
    {
        DB::enableQueryLog();
        config(['google_sqlcommenter.include.route' => false]);
        $response = $this->get('api/db/select');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~select 1/\*framework='laravel-\d*.\d*.\d*',controller='RawDBController',action='select',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }
}
