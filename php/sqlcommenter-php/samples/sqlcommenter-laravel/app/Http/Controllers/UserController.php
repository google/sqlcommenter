<?php

namespace App\Http\Controllers;

use App\Http\Controllers\Controller;
use App\Models\User;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class UserController extends Controller
{
    /**
     * Store a new flight in the database.
     *
     * @param  \Illuminate\Http\Request  $request
     * @return \Illuminate\Http\Response
     */
    public function select()
    {
        User::all();
        return 'Success';
    }

    public function create()
    {
        $user_update = new User;
        $user_update->name = 'john';
        $user_update->email ='contact_me@d.com';
        $user_update->password ='Password$3456';
        $user_update->save();
        return 'Success';
    }

    public function delete()
    {
        User::where('email', 'contact_me@d.com')->delete();
        return 'Success';
    }

    public function update()
    {
        User::where('email', 'contact_me@d.com')->update(["name" =>'johnny']);
        return 'Success';
    }
}
