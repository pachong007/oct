<?php
/**
 * Created by PhpStorm.
 * User: night
 * Date: 2021/5/31
 * Time: 16:53
 */

namespace App\Models\Manga;

use Illuminate\Database\Eloquent\Model;

class MangaChapter extends Model
{
    protected $connection = 'manhua';
    protected $table = "bw_comic_chapter_";
    public $timestamps = false;

    public function setTable($tab)
    {
        $this->table = $this->table.$tab;
        return $this;
    }
}
