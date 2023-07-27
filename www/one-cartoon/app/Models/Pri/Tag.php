<?php
/**
 * Created by PhpStorm.
 * User: night
 * Date: 2021/5/31
 * Time: 16:53
 */

namespace App\Models\Pri;

use App\Models\BaseModel;

class Tag extends BaseModel
{
    protected $connection = 'hhlz_private';
    protected $table = "bw_comic_tag";
}
