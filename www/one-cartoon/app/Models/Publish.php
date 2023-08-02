<?php
/**
 * Created by PhpStorm.
 * User: night
 * Date: 2021/5/31
 * Time: 16:53
 */

namespace App\Models;


use Illuminate\Database\Eloquent\Model;

class Publish extends Model
{
    protected $table = 'publish';

    public function getChapterIdAttribute($v)
    {
        return (array)json_decode($v);
    }

    public function getPublishChapterIdAttribute($v)
    {
        return (array)json_decode($v);
    }
}
