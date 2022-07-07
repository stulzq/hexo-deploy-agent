# compile for version
set -e
#make
#if [ $? -ne 0 ]; then
#    echo "make error"
#    exit 1
#fi

app_version=`cat ./VERSION.txt`
bin_prefix="hexo_deploy_agent"
echo "build version: $app_version"

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin'
arch_all='amd64 arm arm64'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        app_dir_name="${bin_prefix}_${app_version}_${os}_${arch}"
        app_path="./packages/${bin_prefix}_${app_version}_${os}_${arch}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./${bin_prefix}_${os}_${arch}.exe" ]; then
                continue
            fi
            mkdir ${app_path}
            mv ./${bin_prefix}_${os}_${arch}.exe ${app_path}/${bin_prefix}.exe
        else
            if [ ! -f "./${bin_prefix}_${os}_${arch}" ]; then
                continue
            fi
            mkdir ${app_path}
            mv ./${bin_prefix}_${os}_${arch} ${app_path}/${bin_prefix}
        fi
        cp ../LICENSE ${app_path}

        mkdir ${app_path}/conf
        cp -rf ../conf/* ${app_path}/conf

        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${app_dir_name}.zip ${app_dir_name}
        else
            tar -zcf ${app_dir_name}.tar.gz ${app_dir_name}
        fi
        cd ..
        rm -rf ${app_path}
    done
done

cd -