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
    public function test_select()
    {
        DB::enableQueryLog();
        $response = $this->get('api/user/select');
        $myDebugVar = DB::getQueryLog();
        fwrite(STDERR, print_r($myDebugVar, TRUE));
        $response->assertStatus(200);
    }
}
