import { request } from "./request.js";

let state;
let isOnline;

// 获取国际化文本
const getMessage = (key) => {
  const i18nContainer = document.getElementById("i18n-messages");
  if (i18nContainer && i18nContainer.dataset[key]) {
    return i18nContainer.dataset[key];
  }
  return "";
};

const showErrorMsg = (msg) => {
  const errorBox = document.getElementById("errorBox");
  const errorMsgContent = document.getElementById("errorMsgContent");
  errorMsgContent.innerHTML = msg;
  errorBox.classList.toggle("active");
  setTimeout(() => {
    errorBox.classList.toggle("active");
  }, 3000);
  setTimeout(() => {
    errorMsgContent.innerHTML = '';
  }, 4000);
};

// 表单验证函数
const validateForm = (ignoreItems = []) => {
  const formItems = document.querySelectorAll(".form-item");
  let isValid = true;

  // 重置所有错误状态
  formItems.forEach((item) => item.classList.remove("error"));

  // 验证必填字段
  document.querySelectorAll(".form-item.required").forEach((formItem) => {
    const input = formItem.querySelector("input");
    if (!input) return;
    const name = input.getAttribute("name");
    if (ignoreItems.includes(name)) return;
    const value = input.value.trim();
    const placeholder = input.getAttribute("placeholder");
    const errDom = document.getElementById(`${name}Error`);

    if (!value) {
      if (name === "email") return;
      formItem.classList.add("error");
      errDom.textContent = placeholder;
      isValid = false;
    } else {
      switch (name) {
        case "phone":
          if (!/^1[3-9]\d{9}$/.test(value)) {
            formItem.classList.add("error");
            errDom.textContent = getMessage("pleaseenterrealphonenumber");
            isValid = false;
          } else {
            errDom.textContent = "";
          }
          break;
        case "email":
          if (!/^[\w-]+(\.[\w-]+)*@([\w-]+\.)+[a-zA-Z]{2,7}$/.test(value)) {
            formItem.classList.add("error");
            errDom.textContent = getMessage("pleaseenterrealemail");
            isValid = false;
          } else {
            errDom.textContent = "";
          }
          break;
        default:
          errDom.textContent = "";
          break;
      }
    }
  });

  return isValid;
};

const getFormValues = (ignoreItems = []) => {
  const data = {};
  document.querySelectorAll(".form-item").forEach((formItem) => {
    const input =
      formItem.querySelector("input") || formItem.querySelector("select");
    const name = input.getAttribute("name");
    if (!name || ignoreItems.includes(name)) return;
    data[name] = input.value;
  });
  return data;
};

const initInstallAuth = async () => {
  const confirmBtn = document.getElementById("submitAuth");

  // 行业下拉框
  // const industrySelect = document.getElementById("industryCode");
  // const verifyCodeFormItem = document.getElementById("verifyCodeFormItem");
  // const optionsData = [
  //   { value: "A", text: "化工" },
  //   { value: "B", text: "石化" },
  //   { value: "C", text: "电力" },
  //   { value: "D", text: "冶金" },
  //   { value: "E", text: "建材" },
  //   { value: "F", text: "造纸" },
  //   { value: "G", text: "食品医药" },
  //   { value: "H", text: "装备制造" },
  //   { value: "I", text: "汽车汽配" },
  //   { value: "J", text: "金属制品" },
  //   { value: "K", text: "电子电器与仪器仪表" },
  //   { value: "N", text: "教学科研" },
  //   { value: "O", text: "智慧园区" },
  //   { value: "P", text: "信息软件生态伙伴" },
  //   { value: "Q", text: "专项课题" },
  //   { value: "R", text: "公用服务" },
  //   { value: "S", text: "城市基建" },
  //   { value: "L", text: "新材料行业" },
  //   { value: "M", text: "其他行业" },
  // ];
  // const fragment = document.createDocumentFragment();
  // optionsData.forEach((item) => {
  //   const option = document.createElement("option");
  //   option.value = item.value;
  //   option.textContent = item.text;
  //   fragment.appendChild(option);
  // });
  // industrySelect.appendChild(fragment);

  // 离线需要认证码
  if (!isOnline) {
    verifyCodeFormItem.style.display = "block";
    // 获取认证码
    document
      .getElementById("getCodeBtn")
      .addEventListener("click", function () {
        if (!validateForm(["verifyCode"])) return;
        const formValues = getFormValues(["verifyCode"]);

        request("/inter-api/supos/license/register/url")
          .then((data) => {
            document.getElementById("codeErrorMsg").style.display = "block";
            document.getElementById("qrImg").style.display = "block";
            document.getElementById("qrImg").innerHTML = "";

            // 处理 URL
            let finalUrl = `${data}&c=${encodeURIComponent(
              formValues.companyName
            )}&u=${encodeURIComponent(
              formValues.username
            )}&p=${encodeURIComponent(formValues.phone)}&e=${encodeURIComponent(
              formValues.email
            )}&i=${encodeURIComponent(formValues.industryCode)}`;

            new QRCode("qrImg", {
              text: finalUrl,
              width: 120,
              height: 120,
              correctLevel: 3,
            });
          })
          .catch((error) => {
            showErrorMsg(error.message);
          });
      });
  } else {
    verifyCodeFormItem.parentNode.removeChild(verifyCodeFormItem);
  }

  // 输入框实时验证
  document.querySelectorAll(".form-item.required input").forEach((input) => {
    input.addEventListener("input", function (e) {
      const formItem = this.closest(".form-item");
      const value = e.target.value.trim();
      const name = e.target.name;
      const placeholder = e.target.placeholder;
      const errDom = document.getElementById(`${name}Error`);
      formItem.classList.remove("error");
      if (!value) {
        if (name === "email") return;
        formItem.classList.add("error");
        errDom.textContent = placeholder;
      } else {
        switch (name) {
          case "phone":
            if (!/^1[3-9]\d{9}$/.test(value)) {
              formItem.classList.add("error");
              errDom.textContent = getMessage("pleaseenterrealphonenumber");
            } else {
              errDom.textContent = "";
            }
            break;
          case "email":
            if (!/^[\w-]+(\.[\w-]+)*@([\w-]+\.)+[a-zA-Z]{2,7}$/.test(value)) {
              formItem.classList.add("error");
              errDom.textContent = getMessage("pleaseenterrealemail");
            } else {
              errDom.textContent = "";
            }
            break;
          default:
            errDom.textContent = "";
            break;
        }
      }
    });
  });
  confirmBtn.addEventListener("click", function () {
    if (!validateForm()) return;
    const data = getFormValues();
    request("/inter-api/supos/license/verify", {
      method: "POST",
      body: JSON.stringify(data),
    })
      .then(() => {
        window.location.reload();
      })
      .catch((error) => {
        showErrorMsg(error.message);
      });
  });
};

const loginContainer = document.querySelector(".pf-v5-c-login__container");
const installAuthBox = document.getElementById("authBox");
const mainBox = document.getElementById("main");
const loadingWrap = document.getElementById("loadingWrap");
loadingWrap.style.display = "flex";

(function() {
  installAuthBox.style.display = "none";
  mainBox.style.display = "block";
  loadingWrap.style.display = "none";
  loadingWrap.parentNode.removeChild(loadingWrap);
})()
// request(`/inter-api/supos/license/state`)
//   .then((data) => {
//     state = data.state;
//     isOnline = data.online;
//
//     // 0-未授权未激活  1-已授权未激活 2-未授权已激活（一般不会有）3-已激活已授权
//     if (state === 1) {
//       installAuthBox.style.display = "block";
//       mainBox.style.display = "none";
//       loginContainer.style.alignItems = "center";
//       loginContainer.style.gridTemplateAreas = '"main"';
//       initInstallAuth();
//     } else {
//       installAuthBox.style.display = "none";
//       mainBox.style.display = "block";
//     }
//     loadingWrap.style.display = "none";
//     loadingWrap.parentNode.removeChild(loadingWrap);
//   })
//   .catch(() => {
//     installAuthBox.style.display = "none";
//     mainBox.style.display = "block";
//     loadingWrap.style.display = "none";
//     loadingWrap.parentNode.removeChild(loadingWrap);
//   });
