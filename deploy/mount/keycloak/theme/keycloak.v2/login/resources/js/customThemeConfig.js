import { request } from "./request.js";

const checkUrl = (dom, url, defaultUrl) => {
  if(url){
    dom.src = url;
    dom.onerror = function () {
      this.onerror = null;
      if (defaultUrl) {
        this.src = defaultUrl;
      } else {
        this.style.display = 'none'
      }
    };
  }else{
    if (defaultUrl) {
      dom.src = defaultUrl;
    } else {
      dom.style.display = 'none'
    }
  }
};

const OriginUrl = '/files/system/resource/supos'

export const handleTheme = async (keycloakUrl, lang) => {
  const DARK_MODE_CLASS = "pf-v5-theme-dark";
  const { classList } = document.documentElement;
  // 背景图
  const backgroundImgDom = document.querySelector(".pf-v5-c-login-left");
  // slogan
  const sloganImgDom = document.querySelector(".supos-logo");
  // logo
  const logoImgDom = document.querySelector(".supos-login-slogan");

  const favicon = document.getElementById("dynamic-favicon");

  function updateDarkMode(isDark, themeConfig) {
    const {
      general,
      dark,
      light,
    } = themeConfig || {};
    // 浏览器图标
    // 检测 SVG 是否加载失败
    const faviconIcon = new Image();
    if (general?.browserIco) {
      favicon.href = `${OriginUrl}/${general.browserIco}`;
    }
    faviconIcon.src = favicon.href;
    faviconIcon.onerror = function () {
      // 替换为备用图标
      favicon.href = "/logo-ico.svg";
    };

    if (isDark) {
      //暗色主题
      classList.add(DARK_MODE_CLASS);
      if (logoImgDom && sloganImgDom) {
        checkUrl(
          logoImgDom,
          `${OriginUrl}/${dark?.logo || 'logo-dark.svg'}`,
          `${keycloakUrl}/img/logo-dark.svg`
        );
        checkUrl(
          sloganImgDom,
          `${OriginUrl}/${dark?.loginSloganImg}`,
        );
        backgroundImgDom.style.backgroundImage = dark?.loginBackgroundImg
          ? `url(${OriginUrl}/${dark?.loginBackgroundImg})`
          : `url(${keycloakUrl}/img/login-background.svg)`;
      }
    } else {
      //亮色主题
      classList.remove(DARK_MODE_CLASS);
      if (logoImgDom && sloganImgDom) {
        checkUrl(
          logoImgDom,
          `${OriginUrl}/${light?.logo || 'logo-dark.svg'}`,
          `${keycloakUrl}/img/logo-light.svg`
        );
        checkUrl(
          sloganImgDom,
          `${OriginUrl}/${light?.loginSloganImg}`,
        );

        backgroundImgDom.style.backgroundImage = light?.loginBackgroundImg
          ? `url(${OriginUrl}/${light?.loginBackgroundImg})`
          : `url(${keycloakUrl}/img/login-background.svg)`;
      }
    }
  }

  const useDefaultTheme = () => {
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    updateDarkMode(mediaQuery.matches);
    mediaQuery.addEventListener("change", (event) =>
      updateDarkMode(event.matches)
    );
  };

  try {
    const themeConfig = await request(`${OriginUrl}/theme-config.json`)
    if (themeConfig) {
      updateDarkMode(themeConfig?.general?.loginPageTheme === 'dark', themeConfig);
    } else {
      useDefaultTheme();
    }
    document.body.style.opacity = 1;
  } catch (err) {
    useDefaultTheme();
    document.body.style.opacity = 1;
    console.log(err);
  }
};
