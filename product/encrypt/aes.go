package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/kataras/iris/v12/x/errors"
)

//高级加密标准，(Advanced Encryption Standard，AES)，几乎不可能破解

// 16,24,32位的密钥长度分别对应128,192,256位的加密强度。
// key不能泄漏
var PwdKey = []byte("DIS**#KKKDJJSKDI")

// PKCS7填充模式
func PKCS7Padding(cipherttext []byte, blockSize int) []byte {
	padding := blockSize - len(cipherttext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherttext, padtext...) //...语法糖，把padtext切片里面的元素逐个追加，而不是把切片整体追加
}

// 填充的反向操作，删除多余的字符串
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误!")
	} else {
		//获取填充字符串长度
		unpadding := int(origData[length-1])
		// fmt.Println("unpadding:", unpadding) //88,length是16,所以报错，这里有问题
		//截取切片，删除填充的字节，并且返回明文
		return origData[:(length - unpadding)], nil
	}
}

// 实现解密
func AesDeCrypt(crypted []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	bolcksize := block.BlockSize()
	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:bolcksize])
	origData := make([]byte, len(crypted))
	//这个函数既可以用来加密也可以用来解密
	blockMode.CryptBlocks(origData, crypted)
	//去除填充字符串
	// fmt.Println("orinData:", origData)
	origData, err = PKCS7UnPadding(origData) //这里有问题
	if err != nil {
		return nil, err
	}
	return origData, nil
}

// 实现加密
func AesEcrypt(origData []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小,16还是32还是64
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = PKCS7Padding(origData, blockSize)
	//采用AES加密方法中CBC加密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	//执行加密
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 加密base64
func EnPwdCode(pwd []byte) (string, error) {
	result, err := AesEcrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

// 解密base64
func DePwdCode(pwd string) ([]byte, error) {
	//解密base64字符串
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	//执行AES解密
	return AesDeCrypt(pwdByte, PwdKey) //这里有问题
}
