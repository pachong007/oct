<?php
/**
 * Created by PhpStorm.
 * User: night
 * Date: 2021/5/31
 * Time: 16:53
 */

namespace App\Models;


class Image extends BaseModel
{
    public $timestamps = false;
    protected $connection = 'mysql';
    protected $table = 'mc_comic_pic';
}
