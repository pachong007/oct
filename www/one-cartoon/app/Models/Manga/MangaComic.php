<?php
/**
 * Created by PhpStorm.
 * User: night
 * Date: 2021/5/31
 * Time: 16:53
 */

namespace App\Models\Manga;

use Illuminate\Database\Eloquent\Model;

class MangaComic extends Model
{
    protected $connection = 'manhua';
    protected $table = "bw_comic";
    protected $primaryKey = 'comic_id';
    public $timestamps = false;
}
