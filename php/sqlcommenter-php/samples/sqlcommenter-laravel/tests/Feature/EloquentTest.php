<?php

namespace Tests\Feature;

use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;
use Illuminate\Support\Facades\DB;

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
        fwrite(STDERR, print_r($myDebugVar, TRUE));
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~delete from `users` where `email` = \?/\*framework='laravel-9.9.0',controller='UserController',action='delete',route='%%2Fapi%%2Fuser%%2Fdelete',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }
    public function test_create()
    {
        DB::enableQueryLog();
        $response = $this->get('api/user/create');
        $myDebugVar = DB::getQueryLog();
        fwrite(STDERR, print_r($myDebugVar, TRUE));
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~insert into `users` \(`name`, `email`, `password`, `updated_at`, `created_at`\) values \(\?, \?, \?, \?, \?\)/\*framework='laravel-9.9.0',controller='UserController',action='create',route='%%2Fapi%%2Fuser%%2Fcreate',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }

    public function test_select_all()
    {
        DB::enableQueryLog();
        $response = $this->get('api/user/select');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~select \* from `users`/\*framework='laravel-9.9.0',controller='UserController',action='select',route='%%2Fapi%%2Fuser%%2Fselect',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }

    public function test_update()
    {
        DB::enableQueryLog();
        $response = $this->get('api/user/update');
        $myDebugVar = DB::getQueryLog();
        $response->assertStatus(200);
        $this->assertMatchesRegularExpression("~update `users` set `name` = \?, `users`.`updated_at` = \? where `email` = \?/\*framework='laravel-9.9.0',controller='UserController',action='update',route='%%2Fapi%%2Fuser%%2Fupdate',db_driver='mysql',traceparent='\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\d{1,2}'\*/~", end($myDebugVar)['query']);
    }
}
