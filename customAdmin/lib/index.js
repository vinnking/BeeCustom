layui.extend({
    setter: 'config' //配置模块
    , admin: 'lib/admin' //核心模块
    , view: 'lib/view' //视图渲染模块
    , tablePlug: 'lib/tablePlug/tablePlug'
    , AutoComplete: 'lib/AutoComplete'
}).define(['setter', 'tablePlug',"admin",'AutoComplete'], function (exports) {
    var setter = layui.setter
        , element = layui.element
        , admin = layui.admin
        , tabsPage = admin.tabsPage
        , view = layui.view

        //打开标签页
        , openTabsPage = function (url, text) {
            //遍历页签选项卡
            var matchTo
                , tabs = $('#LAY_app_tabsheader>li')
                , path = url.replace(/(^http(s*):)|(\?[\s\S]*$)/g, '');

            tabs.each(function (index) {
                var li = $(this)
                    , layid = li.attr('lay-id');

                if (layid === url) {
                    matchTo = true;
                    tabsPage.index = index;
                }
            });

            text = text || '新标签页';

            if (setter.pageTabs) {
                //如果未在选项卡中匹配到，则追加选项卡
                if (!matchTo) {
                    $(APP_BODY).append([
                        '<div class="layadmin-tabsbody-item layui-show">'
                        , '<iframe src="' + url + '" frameborder="0" class="layadmin-iframe" id="' + text + 'iframe"></iframe>'
                        , '</div>'
                    ].join(''));
                    tabsPage.index = tabs.length;
                    element.tabAdd(FILTER_TAB_TBAS, {
                        title: '<span>' + text + '</span>'
                        , id: url
                        , attr: path
                    });
                }
            } else {
                var iframe = admin.tabsBody(admin.tabsPage.index).find('.layadmin-iframe');
                iframe[0].contentWindow.location.href = url;
            }

            //定位当前tabs
            element.tabChange(FILTER_TAB_TBAS, url);
            admin.tabsBodyChange(tabsPage.index, {
                url: url
                , text: text
            });
        }

        , APP_BODY = '#LAY_app_body', FILTER_TAB_TBAS = 'layadmin-layout-tabs'
        , $ = layui.$, $win = $(window);

    //初始
    if (admin.screen() < 2) admin.sideFlexible();

    //ajax 请求头部
    $.ajaxSetup({
        headers: {
            'X-Xsrftoken': $('meta[name=_xsrf]').attr('content')
        }
    });

    //将模块根路径设置为 controller 目录
    layui.config({
        base: setter.base + 'modules/'
    });

    //扩展 lib 目录下的其它模块
    layui.each(setter.extend, function (index, item) {
        var mods = {};
        mods[item] = '{/}' + setter.base + 'lib/extend/' + item;
        layui.extend(mods);
    });

    view().autoRender();


    //复核修改错误内容改变颜色
    $("body").on("change", ".is_check_fail", function(){
        console.log(1);
        $(this).addClass("is_check_ways").removeClass("is_check_fail");
    });

    //给所有input添加清空图标
    $("input[type='text']").each(function(){
        $(this).after(`<i class="layui-icon layui-icon-close caller-dump-icon"></i>`);
    });

    //input获取焦点显示清空图标
    $("body").on("focus", "input[type='text']", function(){
        $(this).next(".caller-dump-icon").show();
        if($(this).siblings(".icon-ic_add").hasClass("icon-ic_add")) {
            $(this).next(".caller-dump-icon").css("right", "50px");
        } else {
            if($(this).siblings(".icon-more-horiz").hasClass("icon-more-horiz")) {
                $(this).next(".caller-dump-icon").css("right", "22px");
            } else {
                $(this).next(".caller-dump-icon").css("right", "2px");
            }
        }
    });
    $("body").on("blur", "input[type='text']", function(){
        setTimeout(() => {
            $(this).next(".caller-dump-icon").hide();
        }, 150);
    });
    $("body").on("click", ".caller-dump-icon", function(){
        $(this).prev().val("");
        $(this).prev().prev().val("");
    });

    //对外输出
    exports('index', {
        openTabsPage: openTabsPage
    });
});