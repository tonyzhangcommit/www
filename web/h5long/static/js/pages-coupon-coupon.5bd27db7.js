(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["pages-coupon-coupon"],{"06d3":function(t,i,n){var o=n("07b5");o.__esModule&&(o=o.default),"string"===typeof o&&(o=[[t.i,o,""]]),o.locals&&(t.exports=o.locals);var a=n("4f06").default;a("642d6446",o,!0,{sourceMap:!1,shadowMode:!1})},"07b5":function(t,i,n){var o=n("24fb");i=o(!1),i.push([t.i,".iconfont[data-v-63972fae]{font-family:iconfont!important;font-size:%?40?%;color:#d32323}.content_box[data-v-63972fae]{height:100%;padding-left:%?30?%;padding-right:%?30?%;padding-top:%?30?%;background-color:#f5f5f5}.banner[data-v-63972fae]{width:100%}.banner_class[data-v-63972fae]{width:100%;height:%?220?%}.coupon_box[data-v-63972fae]{margin-top:%?30?%}.coupon_item[data-v-63972fae]{display:flex;\n\t/* justify-content: space-between; */border-radius:%?20?%;padding-top:%?20?%;padding-bottom:%?20?%}.left_part_box[data-v-63972fae]{background-color:#991a1a;border-radius:%?20?%;width:76%}.left_part[data-v-63972fae]{background-color:#fff;margin-left:%?10?%;border-top-right-radius:%?20?%;border-bottom-right-radius:%?20?%;border-right:%?2?% dashed #dadada;padding-top:%?20?%;padding-bottom:%?20?%;padding-left:%?20?%;font-weight:700;font-size:%?30?%}.title[data-v-63972fae]{color:#000;padding-top:%?6?%;padding-bottom:%?10?%}.condition[data-v-63972fae]{padding-top:%?6?%;padding-bottom:%?6?%;font-size:%?26?%;color:#919191}.right_part[data-v-63972fae]{border-radius:%?20?%;background-color:#fff;flex-grow:1;display:flex;align-items:center;flex-direction:column;justify-content:space-evenly}.money[data-v-63972fae]{color:#d32323}.money_info[data-v-63972fae]{font-size:%?50?%;font-weight:700}.used_it[data-v-63972fae]{font-size:%?26?%;color:#fff;border-radius:%?40?%;background-color:#d32323;padding:%?10?% %?16?% %?10?% %?16?%}.notification[data-v-63972fae]{position:absolute;width:80%;height:60%;top:20%;left:10%;margin-top:-10%;border-radius:%?20?%;background-color:#f3f3f3}.notictext[data-v-63972fae]{margin-top:10%;padding-right:%?30?%;padding-left:%?30?%}.icon_type[data-v-63972fae]{height:%?50?%;padding-top:%?20?%;padding-bottom:%?20?%;padding-left:%?20?%;padding-right:%?20?%;margin-left:80%}",""]),t.exports=i},"5a3d":function(t,i,n){"use strict";var o=n("06d3"),a=n.n(o);a.a},"5b2e":function(t,i,n){"use strict";n("7a82"),Object.defineProperty(i,"__esModule",{value:!0}),i.default=void 0;var o={data:function(){return{userinfo:{},is_login:!1,firstopen:!1,coupon_list:[{title:"话费充值满100元减5元",condition:"满100元可用",date:"有效期：2023-11-01至2039-12-31",money:"5",count:"共20张"},{title:"话费充值满200元减10元",condition:"满200元可用",date:"有效期：2023-11-01至2039-12-31",money:"10",count:"共10张"}]}},methods:{closeTip:function(){this.firstopen=!1},usecoupon:function(){this.gotoCardShop()},gotoCardShop:function(){window.location.href="http://my.tangjiuhuichong.com/mobile/pages/huafei/huafei?select=0&item_id=0&invite_code=T4FFP6"}},onLoad:function(){},onShow:function(){}};i.default=o},7557:function(t,i,n){"use strict";n.d(i,"b",(function(){return o})),n.d(i,"c",(function(){return a})),n.d(i,"a",(function(){}));var o=function(){var t=this,i=t.$createElement,o=t._self._c||i;return o("v-uni-view",{staticClass:"content_box"},[o("v-uni-view",{staticClass:"banner"},[o("v-uni-image",{staticClass:"banner_class",attrs:{src:n("f37d")}})],1),o("v-uni-view",{staticClass:"coupon_box"},t._l(t.coupon_list,(function(i,n){return o("v-uni-view",{key:n,staticClass:"coupon_item"},[o("v-uni-view",{staticClass:"left_part_box"},[o("v-uni-view",{staticClass:"left_part"},[o("v-uni-view",{staticClass:"title"},[o("v-uni-text",[t._v(t._s(i.title)),o("v-uni-text",{staticStyle:{"font-size":"20rpx","margin-left":"10rpx"}},[t._v(t._s(i.count))])],1)],1),o("v-uni-view",{staticClass:"condition"},[o("v-uni-text",[t._v(t._s(i.condition))])],1),o("v-uni-view",{staticClass:"condition"},[o("v-uni-text",[t._v(t._s(i.date))])],1)],1)],1),o("v-uni-view",{staticClass:"right_part"},[o("v-uni-view",{staticClass:"money"},[o("v-uni-text",{staticClass:"iconfont icon-renminbi2"}),o("v-uni-text",{staticClass:"money_info"},[t._v(t._s(i.price))])],1),o("v-uni-view",{staticClass:"used_it",on:{click:function(i){arguments[0]=i=t.$handleEvent(i),t.usecoupon()}}},[o("v-uni-text",[t._v("立即使用")])],1)],1)],1)})),1),o("v-uni-view",{staticClass:"bottem_box",staticStyle:{position:"absolute",bottom:"20rpx","font-size":"30rpx"}},[o("v-uni-view",{staticClass:"item_box"},[o("v-uni-view",{staticClass:"title",staticStyle:{"font-size":"26rpx"}},[o("v-uni-view",[t._v("注：请输入正确手机号，点击本人充值时，手机号自动填充为登录手机号")])],1)],1),o("v-uni-view",{staticClass:"item_box"},[o("v-uni-view",{staticClass:"title",staticStyle:{color:"black","font-weight":"bold","font-size":"26rpx"}},[o("v-uni-view",[t._v("客服电话：4009265005 在线时间：9:00-18:00")]),o("v-uni-view",[t._v('其他时间请联系在线客服：点击"我的-在线客服"')])],1)],1),o("v-uni-view",{staticClass:"item_box"},[o("v-uni-view",{staticClass:"title",staticStyle:{color:"black","font-weight":"bold","font-size":"26rpx"}},[o("v-uni-view",[t._v("支付成功后，预计在72小时到账，请关注余额情况")])],1)],1)],1)],1)},a=[]},d159:function(t,i,n){"use strict";n.r(i);var o=n("5b2e"),a=n.n(o);for(var e in o)["default"].indexOf(e)<0&&function(t){n.d(i,t,(function(){return o[t]}))}(e);i["default"]=a.a},d446:function(t,i,n){"use strict";n.r(i);var o=n("7557"),a=n("d159");for(var e in a)["default"].indexOf(e)<0&&function(t){n.d(i,t,(function(){return a[t]}))}(e);n("5a3d");var s=n("f0c5"),d=Object(s["a"])(a["default"],o["b"],o["c"],!1,null,"63972fae",null,!1,o["a"],void 0);i["default"]=d.exports},f37d:function(t,i,n){t.exports=n.p+"static/img/huafeibanner.d3e49fd7.png"}}]);