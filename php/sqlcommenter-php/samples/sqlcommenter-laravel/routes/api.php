<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;
use App\Http\Controllers\UserController;
use App\Http\Controllers\RawDBController;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| is assigned the "api" middleware group. Enjoy building your API!
|
*/

Route::middleware('auth:sanctum')->get('/user', function (Request $request) {
    return $request->user();
});


Route::get('/user/select', [UserController::class, 'select']);
Route::get('/user/create', [UserController::class, 'create']);
Route::get('/user/delete', [UserController::class, 'delete']);
Route::get('/user/update', [UserController::class, 'update']);

Route::get('/db/select', [RawDBController::class, 'select']);
Route::get('/db/insert', [RawDBController::class, 'insert']);
Route::get('/db/delete', [RawDBController::class, 'delete']);
Route::get('/db/update', [RawDBController::class, 'update']);
Route::get('/db/selectone', [RawDBController::class, 'selectOne']);
