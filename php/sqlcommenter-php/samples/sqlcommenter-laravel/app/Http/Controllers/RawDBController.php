<?php

namespace App\Http\Controllers;

use Illuminate\Support\Facades\DB;

class RawDBController extends Controller
{
    public function select(){
        DB::select("select 1;");
    }
    public function insert(){
        DB::insert("insert into users (`name`, `email`, `password`) values ('john', 'contact_me@daa.com', 'Passsword3456');");
    }
    public function update(){
        DB::update("update users set name='johnny' where name='john';");
    }
    public function delete(){
        DB::delete("delete from users where name='johnny';");
    }
    public function selectOne(){
        DB::selectOne("select 1;");
    }
}
