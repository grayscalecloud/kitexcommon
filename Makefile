

common:
	hz model -out_dir consts -model_dir ./ -idl ./idl/enums.thrift



errno:
	hz model -out_dir consts -model_dir ./ -idl ./idl/errno.thrift


test:
	rm -rf consts/test
	hz model -out_dir consts -model_dir ./ -idl ./idl/test.thrift