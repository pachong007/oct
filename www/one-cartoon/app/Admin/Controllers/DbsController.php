<?php

namespace App\Admin\Controllers;

use App\Models\Dbs;
use DLP\Tool\Assistant;
use DLP\Widget\Plane;
use Encore\Admin\Controllers\AdminController;
use Encore\Admin\Form;
use Encore\Admin\Grid;
use Encore\Admin\Layout\Content;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Route;
use Illuminate\Database\Eloquent\Builder;

class DbsController extends AdminController
{
    protected $title = '推送库管理';

    protected function grid()
    {
        $grid = new Grid(new Dbs());


        $grid->column('id', __('ID'))->sortable();
        $grid->column('name', '用户名称');

        /*配置*/
        $grid->disableExport();
        $grid->disableRowSelector();
        /*查询匹配*/
        $grid->filter(function ($filter) {
            $filter->like('name', '名称');
        });
        /*弹窗配置*/
        $url = rtrim(config('app.url'), '/') . '/' . Route::current()->uri;
        $grid->actions(function ($actions) use ($url) {
            $actions->disableView();
            $actions->disableEdit();
        });
        return $grid;
    }



    /**
     * Make a form builder.
     *
     * @return Form
     */
    protected function form($id = '')
    {
        $form = new Form(new Dbs());
        $form->input('name','mcc库名');
        /*配置*/
        CommonController::disableDetailConf($form);
        return $form;
    }
}
