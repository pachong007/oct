<?php
/**
 * Created by PhpStorm.
 * User: night
 * Date: 2021/5/31
 * Time: 16:53
 */

namespace App\Models;


use Illuminate\Support\Facades\DB;

class Comic extends BaseModel
{
    protected $connection = 'mysql';
    protected $table = 'comic';

    public function setDynamicConnection($databaseName, $databaseUsername, $databasePassword)
    {
        DB::connection($this->connection)->disconnect();
        // 切换到动态数据库连接
        $this->setConnection('dynamic', [
            'driver' => 'mysql',
            'host' => 'localhost',
            'port' => '3306',
            'database' => $databaseName,
            'username' => $databaseUsername,
            'password' => $databasePassword,
        ]);

        DB::reconnect($this->connection);
    }
}
