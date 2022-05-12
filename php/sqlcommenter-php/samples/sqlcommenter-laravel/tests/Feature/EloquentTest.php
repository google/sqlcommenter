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

class EloquentTest extends TestCase
{
    /**
     * A basic test example.
     *
     * @return void
     */
    public function test_the_application_returns_a_successful_response()
    {

        $response = $this->get('/');
        $response->assertStatus(200);
    }

    public function test_delete()
    {
        DB::enableQueryLog();
        $response = $this->get('api/user/delete');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~delete from `users` where `email` = \?/\*framework='laravel-\d*.\d*.\d*',controller='UserController',action='delete',route='%%2Fapi%%2Fuser%%2Fdelete',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }
    public function test_create()
    {
        DB::enableQueryLog();
        $response = $this->get('api/user/create');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~insert into `users` \(`name`, `email`, `password`, `updated_at`, `created_at`\) values \(\?, \?, \?, \?, \?\)/\*framework='laravel-\d*.\d*.\d*',controller='UserController',action='create',route='%%2Fapi%%2Fuser%%2Fcreate',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }

    public function test_select_all()
    {
        DB::enableQueryLog();
        $response = $this->get('api/user/select');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~select \* from `users`/\*framework='laravel-\d*.\d*.\d*',controller='UserController',action='select',route='%%2Fapi%%2Fuser%%2Fselect',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }

    public function test_update()
    {
        DB::enableQueryLog();
        $response = $this->get('api/user/update');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~update `users` set `name` = \?, `users`.`updated_at` = \? where `email` = \?/\*framework='laravel-\d*.\d*.\d*',controller='UserController',action='update',route='%%2Fapi%%2Fuser%%2Fupdate',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }

    public function test_select_all_db_driver_disabled()
    {
        DB::enableQueryLog();
        config(['google_sqlcommenter.include.db_driver' => false]);
        $response = $this->get('api/user/select');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~select \* from `users`/\*framework='laravel-\d*.\d*.\d*',controller='UserController',action='select',route='%%2Fapi%%2Fuser%%2Fselect',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }
}
