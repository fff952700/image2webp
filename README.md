### 从db中获取image url 下载到内存中格式化为webp后上传到S3
1、通过filter中获取db的表明和where条件。返回interface写死id,code,image,interface自行修改  
2、filter中设置了并发数和split分割code的条件并且返回最后一个  
3、使用了chai2010 webp 需要安装gcc解析器  
4、format格式通过switch case定义。目前只定义了png,jpg,jpeg,webp  
5、配置文件使用的config.toml放在根目录